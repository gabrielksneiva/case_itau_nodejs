import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';


@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatButtonModule,
    MatIconModule
  ],
  template: `
    <header class="app-header">
      <h1>Itau Case Manager</h1>
      <button 
        mat-icon-button
        (click)="toggleTheme()"
        aria-label="Toggle dark/light theme"
        class="theme-toggle"
      >
        <mat-icon class="theme-icon" [ngClass]="{'dark': isDarkMode}">{{ isDarkMode ? 'light_mode' : 'dark_mode' }}</mat-icon>
      </button>
    </header>

    <main>
      <router-outlet></router-outlet>
    </main>
  `,
  styles: [`
    .app-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1rem 2rem;
      background-color: var(--bg-card);
      color: var(--text-primary);
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
      position: sticky;
      top: 0;
      z-index: 1000;
      height: 64px; /* Opcional, mas ajuda a definir uma altura fixa */
    }

    /* 👇 ALTERAÇÃO AQUI 👇 */
    main {
      padding: 2rem;
      padding-top: calc(64px + 2rem); /* Altura do header + espaçamento desejado */
    }
  `]
})
export class AppComponent {
  isDarkMode = false;

  constructor() {
    // Definindo o tema inicial como 'light' para consistência com isDarkMode = false
    document.documentElement.setAttribute('data-theme', 'light');
  }

  toggleTheme() {
    this.isDarkMode = !this.isDarkMode;
    document.documentElement.setAttribute('data-theme', this.isDarkMode ? 'dark' : 'light');
  }
}