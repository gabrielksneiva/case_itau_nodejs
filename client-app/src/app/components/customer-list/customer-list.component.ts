import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { CustomerService, Customer } from '../../services/customer.service';
import { ConfirmationDialogComponent } from '../confirmation-dialog/confirmation-dialog.component';

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
    private customerService: CustomerService,
    public dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.loadCustomers();
  }

  loadCustomers() {
    this.customerService.getAll().subscribe({
      next: (list: Customer[]) => {
        // inicializa campos de UI que podem não existir vindo do backend
        this.customers = list.map(c => ({
          ...c,
          saldoOculto: c.saldoOculto ?? false,
          // se o backend já enviar saldoUpdatedAt, preserva; senão, undefined
          saldoUpdatedAt: (c as any).saldoUpdatedAt ?? undefined
        }));
      },
      error: (err) => {
        console.error('Erro ao carregar clientes', err);
      }
    });
  }

  openDeleteDialog(customerId: number, customerName: string): void {
    const dialogRef = this.dialog.open(ConfirmationDialogComponent, {
      width: '450px',
      data: {
        title: 'Confirmar Exclusão',
        message: `Você tem certeza que deseja excluir o cliente "${customerName}"?`
      }
    });
    
    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.customerService.delete(customerId).subscribe(() => {
          this.loadCustomers();
        });
      }
    });
  }

   trackByCustomer(index: number, c: Customer): string | number {
    return c?.id ?? index;
  }
}

