import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, Validators, FormGroup } from '@angular/forms';
import { CustomerService, Customer } from '../../services/customer.service';
import { Router, ActivatedRoute } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressBarModule } from '@angular/material/progress-bar';

@Component({
  selector: 'app-customer-form',
  standalone: true,
  imports: [CommonModule,
    ReactiveFormsModule,
    MatCardModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatProgressBarModule
  ],
  templateUrl: './customer-form.component.html'
})
export class CustomerFormComponent implements OnInit {
  form!: FormGroup; // inicializamos depois no constructor
  isEdit = false;
  id?: number;
  loading = false;
  error?: string;

  constructor(
    private fb: FormBuilder,
    private service: CustomerService,
    private router: Router,
    private route: ActivatedRoute
  ) {
    // Inicializa o formulário aqui, após FormBuilder já ter sido injetado
    this.form = this.fb.group({
      nome: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
    });
  }

  ngOnInit(): void {
    const idStr = this.route.snapshot.paramMap.get('id');
    if (idStr) {
      this.isEdit = true;
      this.id = Number(idStr);
      this.loading = true;
      this.service.getById(this.id).subscribe({
        next: (c) => {
          this.form.patchValue({ nome: c.nome, email: c.email, saldo: c.saldo });
          this.loading = false;
        },
        error: (err) => {
          this.error = 'Erro ao carregar cliente';
          this.loading = false;
        }
      });
    }
  }

  submit() {
    if (this.form.invalid) return;
    this.loading = true;
    const payload = this.form.value;
    const obs = this.isEdit && this.id
      ? this.service.update(this.id, payload)
      : this.service.create(payload);

    obs.subscribe({
      next: () => this.router.navigate(['/clientes']),
      error: (err) => {
        this.error = 'Erro ao salvar cliente';
        this.loading = false;
      }
    });
  }

  cancel() {
    this.router.navigate(['/clientes']);
  }
}
