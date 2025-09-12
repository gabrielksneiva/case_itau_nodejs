import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Customer {
  id: number;
  nome: string;
  email: string;
  saldo: number;
  saldoOculto?: boolean;                 
  saldoUpdatedAt?: string | Date;  
}

@Injectable({ providedIn: 'root' })
export class CustomerService {
  private apiUrl = 'http://localhost:8080/clientes';

  constructor(private http: HttpClient) {}

  getAll(): Observable<Customer[]> {
    return this.http.get<Customer[]>(this.apiUrl);
  }

  getById(id: number): Observable<Customer> {
    return this.http.get<Customer>(`${this.apiUrl}/${id}`);
  }

  create(payload: Partial<Customer>): Observable<Customer> {
    return this.http.post<Customer>(this.apiUrl, payload);
  }

  update(id: number, payload: Partial<Customer>): Observable<Customer> {
    return this.http.put<Customer>(`${this.apiUrl}/${id}`, payload);
  }

  delete(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }

  deposit(id: number, valor: number): Observable<Customer> {
    return this.http.post<Customer>(`${this.apiUrl}/${id}/depositar`, { valor });
  }

  withdraw(id: number, valor: number): Observable<Customer> {
    return this.http.post<Customer>(`${this.apiUrl}/${id}/sacar`, { valor });
  }
}
