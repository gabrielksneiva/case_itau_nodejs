import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatIconModule } from '@angular/material/icon';

export interface TransactionConfig {
  title: string;
  buttonText: string;
  color: 'primary' | 'accent' | 'warn';
  icon: string;
}

@Component({
  selector: 'app-transaction-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatProgressBarModule,
    MatIconModule
  ],
  templateUrl: './transaction-form.component.html'
})
export class TransactionFormComponent {
  @Input() config!: TransactionConfig;
  @Input() clientName!: string;
  @Input() loading = false;
  @Input() error?: string;

  @Output() formSubmit = new EventEmitter<{ amount: number }>();
  @Output() cancel = new EventEmitter<void>();

  form: FormGroup;

  constructor(private fb: FormBuilder) {
    this.form = this.fb.group({
      amount: [null, [Validators.required, Validators.min(0.01)]]
    });
  }

  submitTransaction() {
    if (this.form.invalid) return;
    this.formSubmit.emit({ amount: this.form.value.amount });
  }

  cancelTransaction() {
    this.cancel.emit();
  }
}

