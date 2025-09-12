import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { TransactionFormComponent, TransactionConfig } from '../transaction-form/transaction-form.component';
import { CustomerService, Customer } from '../../services/customer.service';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';

@Component({
  selector: 'app-customer-deposit',
  standalone: true,
  imports: [
    CommonModule,
    TransactionFormComponent,
    MatSnackBarModule
  ],
  templateUrl: './customer-deposit.component.html'
})
export class CustomerDepositComponent implements OnInit {
  customer: Customer | undefined;
  isLoading = false;
  errorMessage?: string;

  depositConfig: TransactionConfig = {
    title: 'Realizar Depósito',
    buttonText: 'Confirmar Depósito',
    color: 'primary',
    icon: 'arrow_downward'
  };

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private customerService: CustomerService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    const idStr = this.route.snapshot.paramMap.get('id');
    if (idStr) {
      const customerId = Number(idStr);
      this.customerService.getById(customerId).subscribe({
        next: data => this.customer = data,
        error: () => this.router.navigate(['/clientes'])
      });
    } else {
      this.router.navigate(['/clientes']);
    }
  }

  handleTransaction(formData: { valor: number }): void {
    if (!this.customer) return;

    this.isLoading = true;
    this.errorMessage = undefined;
    const valor = formData.valor;
    
    this.customerService.deposit(this.customer.id, valor).subscribe({
      next: () => {
        this.snackBar.open('Depósito realizado com sucesso', 'Fechar', {
          duration: 4000,
          horizontalPosition: 'center',
          verticalPosition: 'top',
          panelClass: ['snackbar-success']
        });
        this.router.navigate(['/clientes']);
      },
      error: (err) => {
        const msg = err.error?.message || err.error?.erro || 'Erro inesperado';
        this.snackBar.open(msg, 'Fechar', {
          duration: 5000,
          horizontalPosition: 'center',
          verticalPosition: 'top',
          panelClass: ['snackbar-error']
        });
        this.isLoading = false;
      }
    });
  }

  goBack(): void {
    this.router.navigate(['/clientes']);
  }
}

