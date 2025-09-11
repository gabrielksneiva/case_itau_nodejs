import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { TransactionFormComponent, TransactionConfig } from '../transaction-form/transaction-form.component';
import { CustomerService, Customer } from '../../services/customer.service';

@Component({
  selector: 'app-customer-deposit',
  standalone: true,
  imports: [
    CommonModule,
    TransactionFormComponent
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
    private customerService: CustomerService
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

  handleTransaction(formData: { amount: number }): void {
    if (!this.customer) return;

    this.isLoading = true;
    this.errorMessage = undefined;
    const amount = formData.amount;
    
    this.customerService.deposit(this.customer.id, amount).subscribe({
      next: () => this.router.navigate(['/clientes']),
      error: (err) => {
        this.errorMessage = err.error?.message || 'Erro ao realizar o depósito. Tente novamente.';
        this.isLoading = false;
      }
    });
  }

  goBack(): void {
    this.router.navigate(['/clientes']);
  }
}

