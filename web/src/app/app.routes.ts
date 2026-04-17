import { Routes } from '@angular/router';
import { InvoicePageComponent } from './pages/invoice-page.component';
import { ProductsPageComponent } from './pages/products-page.component';

export const routes: Routes = [
  { path: '', pathMatch: 'full', redirectTo: 'produtos' },
  { path: 'produtos', component: ProductsPageComponent },
  { path: 'notas-fiscais', component: InvoicePageComponent }
];
