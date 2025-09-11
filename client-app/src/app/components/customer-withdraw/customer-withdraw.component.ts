import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { TransactionFormComponent, TransactionConfig } from '../transaction-form/transaction-form.component';
import { CustomerService, Customer } from '../../services/customer.service';

@Component({
  selector: 'app-customer-withdraw',
  standalone: true,
  imports: [
    CommonModule,
    TransactionFormComponent
  ],
  templateUrl: './customer-withdraw.component.html'
})
export class CustomerWithdrawComponent implements OnInit {
  customer: Customer | undefined;
  isLoading = false;
  errorMessage?: string;

  withdrawConfig: TransactionConfig = {
    title: 'Realizar Saque',
    buttonText: 'Confirmar Saque',
    color: 'warn',
    icon: 'arrow_upward'
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
    
    this.customerService.withdraw(this.customer.id, amount).subscribe({
      next: () => this.router.navigate(['/clientes']),
      error: (err) => {
        this.errorMessage = err.error?.message || 'Saldo insuficiente ou erro na operação. Tente novamente.';
        this.isLoading = false;
      }
    });
  }

  goBack(): void {
    this.router.navigate(['/clientes']);
  }
}

