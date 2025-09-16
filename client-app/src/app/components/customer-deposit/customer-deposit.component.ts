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
    const id = this.route.snapshot.paramMap.get('id');
    if (!id) {
      this.router.navigate(['/clientes']);
      return;
    }

    this.customerService.getById(id).subscribe({
      next: data => this.customer = data,
      error: () => this.router.navigate(['/clientes'])
    });
  }

  handleTransaction(formData: { amount: number }) {
    if (!this.customer) return;
    this.isLoading = true;
    this.errorMessage = undefined;

    this.customerService.deposit(this.customer.id, formData.amount).subscribe({
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
        this.errorMessage = err.error?.message || 'Erro inesperado';
        this.isLoading = false;
      }
    });
  }

  goBack() {
    this.router.navigate(['/clientes']);
  }
}
