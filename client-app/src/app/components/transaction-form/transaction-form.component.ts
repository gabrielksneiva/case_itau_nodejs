import { Component, Input, Output, EventEmitter, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';

// Importações dos módulos do Angular Material necessários para o template
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressBarModule } from '@angular/material/progress-bar';

// Interface para definir a estrutura do objeto de configuração
// Isso garante segurança de tipo e autocompletar ao usar o componente
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


  @Output() formSubmit = new EventEmitter<{ amount: number }>();

  @Output() cancel = new EventEmitter<void>();

  form!: FormGroup;

  constructor(private fb: FormBuilder) {}

  ngOnInit(): void {
    // Inicializa o formulário quando o componente é criado
    this.form = this.fb.group({
      // O valor inicial é null para um campo limpo
      // Validadores: o campo é obrigatório e o valor mínimo deve ser 0.01
      amount: [null, [Validators.required, Validators.min(0.01)]]
    });
  }

  /**
   * Método chamado quando o formulário é submetido.
   * Valida o formulário e, se válido, emite o evento 'formSubmit'.
   */
  submitTransaction(): void {
    // Impede a submissão se o formulário estiver inválido
    if (this.form.invalid) {
      // O Angular Material já exibe os erros, mas você poderia adicionar lógica extra aqui
      return;
    }

    // Emite o valor do formulário para o componente pai, que cuidará da lógica de negócio
    this.formSubmit.emit(this.form.value);
  }
}
