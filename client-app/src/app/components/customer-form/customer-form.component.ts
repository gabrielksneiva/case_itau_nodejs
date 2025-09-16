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
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';

@Component({
  selector: 'app-customer-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatCardModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatProgressBarModule,
    MatSnackBarModule
  ],
  templateUrl: './customer-form.component.html'
})
export class CustomerFormComponent implements OnInit {
  form!: FormGroup;
  isEdit = false;
  id?: string;
  loading = false;
  error?: string;

  constructor(
    private fb: FormBuilder,
    private service: CustomerService,
    private router: Router,
    private route: ActivatedRoute,
    private snackBar: MatSnackBar
  ) {
    this.form = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
    });
  }

  ngOnInit(): void {
    const idStr = this.route.snapshot.paramMap.get('id');
    if (idStr) {
      this.isEdit = true;
      this.id = idStr;
      this.loading = true;
      this.service.getById(this.id).subscribe({
        next: (c: Customer) => {
          this.form.patchValue({ name: c.name, email: c.email });
          this.loading = false;
        },
        error: () => {
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
      next: () => {
        const msg = this.isEdit
          ? 'Cliente atualizado com sucesso'
          : 'Cliente criado com sucesso';

        this.snackBar.open(msg, 'Fechar', {
          duration: 4000,
          horizontalPosition: 'center',
          verticalPosition: 'top',
          panelClass: ['snackbar-success']
        });

        this.router.navigate(['/clientes']);
      },
      error: (err) => {
        const msg = err.error?.message || 'Erro inesperado';
        this.snackBar.open(msg, 'Fechar', {
          duration: 5000,
          horizontalPosition: 'center',
          verticalPosition: 'top',
          panelClass: ['snackbar-error']
        });
        this.loading = false;
      }
    });
  }

  cancel() {
    this.router.navigate(['/clientes']);
  }

  goToTransactions() {
    if (this.id) {
      this.router.navigate([`/clientes/${this.id}/transacoes`]);
    }
  }
}
