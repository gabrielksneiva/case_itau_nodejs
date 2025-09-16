import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { CustomerService, Customer } from '../../services/customer.service';
import { ConfirmationDialogComponent } from '../confirmation-dialog/confirmation-dialog.component';
import { take } from 'rxjs/operators';

@Component({
  selector: 'app-customer-list',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatCardModule,
    MatIconModule,
    MatButtonModule,
    MatDialogModule
  ],
  templateUrl: './customer-list.component.html'
})
export class CustomerListComponent implements OnInit {
  customers: Customer[] = [];

  constructor(
    private readonly customerService: CustomerService,
    private readonly dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.loadCustomers();
  }

  private loadCustomers(): void {
    this.customerService.getAll().pipe(take(1)).subscribe({
      next: (list: Customer[]) => {
        this.customers = list.map(c => ({
          ...c,
          balanceOculto: c.balanceOculto ?? false,
          balanceUpdatedAt: c.balanceUpdatedAt ?? undefined
        }));
      },
      error: (err) => {
        console.error('Erro ao carregar clientes:', err);
      }
    });
  }

  openDeleteDialog(customerId: string, customerName: string): void {
    this.dialog.open(ConfirmationDialogComponent, {
      width: '450px',
      data: {
        title: 'Confirmar Exclusão',
        message: `Você tem certeza que deseja excluir o cliente "${customerName}"?`
      }
    })
    .afterClosed()
    .pipe(take(1))
    .subscribe(result => {
      if (result) {
        this.customerService.delete(customerId).pipe(take(1)).subscribe({
          next: () => this.loadCustomers(),
          error: (err) => console.error('Erro ao excluir cliente:', err)
        });
      }
    });
  }

  trackByCustomer(index: number, customer: Customer): string {
    return customer.id;
  }
}
