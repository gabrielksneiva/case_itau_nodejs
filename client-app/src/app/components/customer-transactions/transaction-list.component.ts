import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatTableModule } from '@angular/material/table';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { CustomerService, Customer } from '../../services/customer.service';

export interface Transaction {
  id: number;
  type: 'deposit' | 'withdraw';
  amount: number;
  createdAt: string;
}

@Component({
  selector: 'app-transaction-list',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatIconModule,
    MatButtonModule,
    MatProgressBarModule,
    MatSnackBarModule,
    MatTableModule
  ],
  templateUrl: './transaction-list.component.html'
})
export class CustomerTransactionListComponent implements OnInit {
  transactions: Transaction[] = [];
  displayedColumns = ['createdAt', 'type', 'amount'];
  customer?: Customer;
  loading = false;
  error?: string;
  id!: string;

  constructor(
    private service: CustomerService,
    private route: ActivatedRoute,
    private router: Router,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    const idStr = this.route.snapshot.paramMap.get('id');
    if (idStr) {
      this.id = idStr;
      this.loadTransactions();
    }
  }

  loadTransactions() {
  this.loading = true;

  this.service.getById(this.id).subscribe({
    next: (customer) => {
      this.customer = customer;

      // Supondo que você tenha um endpoint separado para buscar transações
      this.service.getTransactions(this.id).subscribe({
        next: (res: any) => {
          this.transactions = res.items.map((tx: any) => ({
            id: tx.transaction_id,
            type: tx.type === 'deposit' ? 'deposit' : 'withdraw',
            amount: Number(tx.value),
            createdAt: tx.created_at,
          }));
          this.loading = false;
        },
        error: () => {
          this.error = 'Erro ao carregar transações';
          this.snackBar.open(this.error, 'Fechar', {
            duration: 4000,
            horizontalPosition: 'center',
            verticalPosition: 'top',
            panelClass: ['snackbar-error']
          });
          this.loading = false;
        }
      });
    },
    error: () => {
      this.error = 'Erro ao carregar cliente';
      this.snackBar.open(this.error, 'Fechar', {
        duration: 4000,
        horizontalPosition: 'center',
        verticalPosition: 'top',
        panelClass: ['snackbar-error']
      });
      this.loading = false;
    }
  });
}


  goBack() {
    this.router.navigate(['/clientes']);
  }
}
