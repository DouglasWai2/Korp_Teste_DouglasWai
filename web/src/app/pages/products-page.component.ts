import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, OnInit, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { firstValueFrom } from 'rxjs';
import { API_URLS } from '../api.config';
import { mapProductError } from '../http-error.utils';

interface Product {
  codigo: string;
  descricao: string;
  saldo: number;
}

@Component({
  selector: 'app-products-page',
  imports: [CommonModule, FormsModule],
  templateUrl: './products-page.component.html',
  styleUrl: './products-page.component.css'
})
export class ProductsPageComponent implements OnInit {
  private readonly http = inject(HttpClient);

  readonly products = signal<Product[]>([]);
  readonly loadingProducts = signal(false);
  readonly submittingProduct = signal(false);
  readonly updatingStockCode = signal<string | null>(null);
  readonly deletingProductCode = signal<string | null>(null);
  readonly feedback = signal('');
  readonly feedbackType = signal<'success' | 'error'>('success');
  readonly stockInputs: Record<string, number> = {};

  readonly productForm = {
    codigo: '',
    descricao: '',
    saldo: 0
  };

  async ngOnInit(): Promise<void> {
    await this.loadProducts();
  }

  async loadProducts(): Promise<void> {
    this.loadingProducts.set(true);
    try {
      const response = await firstValueFrom(
        this.http.get<{ data: Product[] }>(`${API_URLS.estoque}/api/products`)
      );
      this.products.set(response.data ?? []);
    } catch (error) {
      this.showFeedback(mapProductError(error, 'Nao foi possivel carregar os produtos.'), 'error');
    } finally {
      this.loadingProducts.set(false);
    }
  }

  async createProduct(): Promise<void> {
    this.submittingProduct.set(true);
    try {
      await firstValueFrom(this.http.post(`${API_URLS.estoque}/api/products`, this.productForm));
      this.productForm.codigo = '';
      this.productForm.descricao = '';
      this.productForm.saldo = 0;
      this.showFeedback('Produto cadastrado com sucesso.', 'success');
      await this.loadProducts();
    } catch (error) {
      this.showFeedback(mapProductError(error, 'Nao foi possivel cadastrar o produto.'), 'error');
    } finally {
      this.submittingProduct.set(false);
    }
  }

  async incrementStock(codigo: string): Promise<void> {
    const quantidade = this.stockInputs[codigo] ?? 0;

    if (quantidade <= 0) {
      this.showFeedback('Informe uma quantidade maior que zero para adicionar ao saldo.', 'error');
      return;
    }

    this.updatingStockCode.set(codigo);
    try {
      await firstValueFrom(
        this.http.patch(`${API_URLS.estoque}/api/products/${codigo}/increment`, {
          quantidade
        })
      );
      this.stockInputs[codigo] = 0;
      this.showFeedback(`Saldo do produto ${codigo} atualizado com sucesso.`, 'success');
      await this.loadProducts();
    } catch (error) {
      this.showFeedback(mapProductError(error, `Nao foi possivel atualizar o saldo do produto ${codigo}.`), 'error');
    } finally {
      this.updatingStockCode.set(null);
    }
  }

  async deleteProduct(codigo: string): Promise<void> {
    this.deletingProductCode.set(codigo);
    try {
      await firstValueFrom(
        this.http.delete(`${API_URLS.estoque}/api/products/${codigo}`)
      );
      delete this.stockInputs[codigo];
      this.showFeedback(`Produto ${codigo} removido com sucesso.`, 'success');
      await this.loadProducts();
    } catch (error) {
      this.showFeedback(mapProductError(error, `Nao foi possivel remover o produto ${codigo}.`), 'error');
    } finally {
      this.deletingProductCode.set(null);
    }
  }

  private showFeedback(message: string, type: 'success' | 'error'): void {
    this.feedback.set(message);
    this.feedbackType.set(type);
  }
}
