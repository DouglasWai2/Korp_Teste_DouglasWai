import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, OnInit, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { firstValueFrom } from 'rxjs';
import { API_URLS } from '../api.config';
import { mapInvoiceCreationError, mapInvoicePrintError, mapProductError } from '../http-error.utils';

interface Product {
  codigo: string;
}

interface NotaFiscal {
  numero: number;
  status: 'Aberta' | 'Fechada';
}

interface NotaFiscalItemForm {
  codigo_produto: string;
  quantidade: number;
}

@Component({
  selector: 'app-invoice-page',
  imports: [CommonModule, FormsModule],
  templateUrl: './invoice-page.component.html',
  styleUrl: './invoice-page.component.css'
})
export class InvoicePageComponent implements OnInit {
  private readonly http = inject(HttpClient);

  readonly products = signal<Product[]>([]);
  readonly notasFiscais = signal<NotaFiscal[]>([]);
  readonly loadingNotas = signal(false);
  readonly submittingNota = signal(false);
  readonly printingNotaNumero = signal<number | null>(null);
  readonly feedback = signal('');
  readonly feedbackType = signal<'success' | 'error'>('success');

  readonly invoiceForm: { itens: NotaFiscalItemForm[] } = {
    itens: [this.newInvoiceItem()]
  };

  async ngOnInit(): Promise<void> {
    await Promise.all([this.loadProducts(), this.loadNotasFiscais()]);
  }

  addInvoiceItem(): void {
    this.invoiceForm.itens.push(this.newInvoiceItem());
  }

  removeInvoiceItem(index: number): void {
    if (this.invoiceForm.itens.length === 1) {
      return;
    }
    this.invoiceForm.itens.splice(index, 1);
  }

  async loadProducts(): Promise<void> {
    try {
      const response = await firstValueFrom(
        this.http.get<{ data: Product[] }>(`${API_URLS.estoque}/api/products`)
      );
      this.products.set(response.data ?? []);
    } catch (error) {
      this.showFeedback(mapProductError(error, 'Nao foi possivel carregar os produtos para selecao.'), 'error');
    }
  }

  async loadNotasFiscais(): Promise<void> {
    this.loadingNotas.set(true);
    try {
      const response = await firstValueFrom(
        this.http.get<{ data: NotaFiscal[] }>(`${API_URLS.faturamento}/api/faturamento/notas-fiscais`)
      );
      this.notasFiscais.set(response.data ?? []);
    } catch (error) {
      this.showFeedback(mapInvoiceCreationError(error, 'Nao foi possivel carregar as notas fiscais.'), 'error');
    } finally {
      this.loadingNotas.set(false);
    }
  }

  async createNotaFiscal(): Promise<void> {
    this.submittingNota.set(true);
    try {
      await firstValueFrom(
        this.http.post(`${API_URLS.faturamento}/api/faturamento/notas-fiscais`, {
          itens: this.invoiceForm.itens
        })
      );
      this.invoiceForm.itens.splice(0, this.invoiceForm.itens.length, this.newInvoiceItem());
      this.showFeedback('Nota fiscal criada com sucesso.', 'success');
      await this.loadNotasFiscais();
    } catch (error) {
      this.showFeedback(mapInvoiceCreationError(error, 'Nao foi possivel cadastrar a nota fiscal.'), 'error');
    } finally {
      this.submittingNota.set(false);
    }
  }

  async printNotaFiscal(numero: number): Promise<void> {
    this.printingNotaNumero.set(numero);
    try {
      await firstValueFrom(
        this.http.patch(`${API_URLS.faturamento}/api/faturamento/notas-fiscais/${numero}/imprimir`, {})
      );
      this.showFeedback(`Nota fiscal #${numero} impressa com sucesso.`, 'success');
      await this.loadNotasFiscais();
    } catch (error) {
      this.showFeedback(mapInvoicePrintError(error, `Nao foi possivel imprimir a nota fiscal #${numero}.`), 'error');
    } finally {
      this.printingNotaNumero.set(null);
    }
  }

  private newInvoiceItem(): NotaFiscalItemForm {
    return {
      codigo_produto: '',
      quantidade: 1
    };
  }

  private showFeedback(message: string, type: 'success' | 'error'): void {
    this.feedback.set(message);
    this.feedbackType.set(type);
  }
}
