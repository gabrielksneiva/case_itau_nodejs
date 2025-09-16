import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CustomerService } from '../../services/customer.service';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-customer-delete',
  standalone: true,
  imports: [CommonModule],
  template: `
    <section *ngIf="id">
      <h2>Deletar Cliente</h2>
      <p>Tem certeza que deseja deletar o cliente com id {{ id }}?</p>
      <button (click)="confirmDelete()">Sim, deletar</button>
      <button (click)="cancel()">Cancelar</button>
      <p *ngIf="error" style="color:red">{{ error }}</p>
    </section>
  `
})
export class CustomerDeleteComponent implements OnInit {
  id?: string;
  error?: string;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private service: CustomerService
  ) {}

  ngOnInit(): void {
    const idStr = this.route.snapshot.paramMap.get('id');
  }

  confirmDelete() {
    if (!this.id) return;
    this.service.delete(this.id).subscribe({
      next: () => this.router.navigate(['/clientes']),
      error: () => this.error = 'Erro ao deletar'
    });
  }

  cancel() {
    this.router.navigate(['/clientes']);
  }
}
