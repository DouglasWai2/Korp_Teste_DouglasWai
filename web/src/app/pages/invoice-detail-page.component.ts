import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, OnInit, inject, signal } from '@angular/core';
import { RouterLink, ActivatedRoute } from '@angular/router';
import { firstValueFrom } from 'rxjs';
import { API_URLS } from '../api.config';
import { mapInvoiceDetailError, mapInvoicePrintError, mapProductError } from '../http-error.utils';

interface Product {
  codigo: string;
  descricao: string;
}

interface NotaFiscalItem {
  codigo_produto: string;
  quantidade: number;
}

interface NotaFiscalDetail {
  numero: number;
  status: 'Aberta' | 'Fechada';
  itens: NotaFiscalItem[];
}

@Component({
  selector: 'app-invoice-detail-page',
  imports: [CommonModule, RouterLink],
  templateUrl: './invoice-detail-page.component.html',
  styleUrl: './invoice-detail-page.component.css'
})
export class InvoiceDetailPageComponent implements OnInit {
  private readonly http = inject(HttpClient);
  private readonly route = inject(ActivatedRoute);

  readonly notaFiscal = signal<NotaFiscalDetail | null>(null);
  readonly productsByCode = signal<Record<string, string>>({});
  readonly loading = signal(false);
  readonly printing = signal(false);
  readonly feedback = signal('');
  readonly feedbackType = signal<'success' | 'error'>('success');

  async ngOnInit(): Promise<void> {
    await Promise.all([this.loadNotaFiscal(), this.loadProducts()]);
  }

  async loadNotaFiscal(): Promise<void> {
    const numero = this.route.snapshot.paramMap.get('numero');
    if (!numero) {
      this.showFeedback('A nota fiscal informada nao foi encontrada.', 'error');
      return;
    }

    this.loading.set(true);
    try {
      const response = await firstValueFrom(
        this.http.get<{ data: NotaFiscalDetail }>(`${API_URLS.faturamento}/api/faturamento/notas-fiscais/${numero}`)
      );
      this.notaFiscal.set(response.data);
    } catch (error) {
      this.showFeedback(mapInvoiceDetailError(error, 'Nao foi possivel carregar a nota fiscal.'), 'error');
    } finally {
      this.loading.set(false);
    }
  }

  async loadProducts(): Promise<void> {
    try {
      const response = await firstValueFrom(
        this.http.get<{ data: Product[] }>(`${API_URLS.estoque}/api/products`)
      );

      const products = response.data ?? [];
      const byCode = products.reduce<Record<string, string>>((accumulator, product) => {
        accumulator[product.codigo] = product.descricao;
        return accumulator;
      }, {});

      this.productsByCode.set(byCode);
    } catch (error) {
      this.showFeedback(mapProductError(error, 'Nao foi possivel carregar os produtos.'), 'error');
    }
  }

  async printNotaFiscal(): Promise<void> {
    const notaFiscal = this.notaFiscal();
    if (!notaFiscal) {
      return;
    }

    this.printing.set(true);
    try {
      await firstValueFrom(
        this.http.patch(`${API_URLS.faturamento}/api/faturamento/notas-fiscais/${notaFiscal.numero}/imprimir`, {})
      );
      this.showFeedback(`Nota fiscal #${notaFiscal.numero} impressa com sucesso.`, 'success');
      await this.loadNotaFiscal();
    } catch (error) {
      this.showFeedback(mapInvoicePrintError(error, `Nao foi possivel imprimir a nota fiscal #${notaFiscal.numero}.`), 'error');
    } finally {
      this.printing.set(false);
    }
  }

  private showFeedback(message: string, type: 'success' | 'error'): void {
    this.feedback.set(message);
    this.feedbackType.set(type);
  }

  productDescription(codigoProduto: string): string {
    return this.productsByCode()[codigoProduto] ?? 'Descricao indisponivel';
  }
}
