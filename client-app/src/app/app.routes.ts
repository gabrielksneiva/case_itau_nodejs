import { Routes } from '@angular/router';
import { CustomerListComponent } from './components/customer-list/customer-list.component';
import { CustomerFormComponent } from './components/customer-form/customer-form.component';
import { CustomerDeleteComponent } from './components/customer-delete/customer-delete.component';
import { CustomerDepositComponent } from './components/customer-deposit/customer-deposit.component';
import { CustomerWithdrawComponent } from './components/customer-withdraw/customer-withdraw.component';
import { CustomerTransactionListComponent } from './components/customer-transactions/transaction-list.component';


export const routes: Routes = [
  { path: '', redirectTo: 'clientes', pathMatch: 'full' },
  { path: 'clientes', component: CustomerListComponent },
  { path: 'clientes/new', component: CustomerFormComponent },
  { path: 'clientes/:id/edit', component: CustomerFormComponent },
  { path: 'clientes/:id/delete', component: CustomerDeleteComponent },
  { path: 'clientes/:id/deposit', component: CustomerDepositComponent },
  { path: 'clientes/:id/withdraw', component: CustomerWithdrawComponent },
  { path: 'clientes/:id/transacoes', component: CustomerTransactionListComponent }
];
