import { Component, Input, Output, EventEmitter, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressBarModule } from '@angular/material/progress-bar';

export interface TransactionConfig {
  title: string;
  buttonText: string;
  color: 'primary' | 'warn';
  icon: string;
}

@Component({
  selector: 'app-transaction-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatCardModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatProgressBarModule
  ],
  templateUrl: './transaction-form.component.html'
})
export class TransactionFormComponent implements OnInit {
  @Input({ required: true }) config!: TransactionConfig;

  @Input() clientName: string | null = '';

  @Input() loading = false;

  @Input() error?: string;

  @Output() formSubmit = new EventEmitter<{ valor: number }>();

  @Output() cancel = new EventEmitter<void>();

  form!: FormGroup;

  constructor(private fb: FormBuilder) {}

  ngOnInit(): void {
    this.form = this.fb.group({
      valor: [null, [Validators.required, Validators.min(0.01)]]
    });
  }

  submitTransaction(): void {
    if (this.form.invalid) {
      return;
    }

    this.formSubmit.emit(this.form.value);
  }
}
