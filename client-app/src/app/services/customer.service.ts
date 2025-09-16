import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Customer {
  id: string;
  name: string;
  email: string;
  balance: number;
  balanceOculto?: boolean;                 
  balanceUpdatedAt?: string | Date;  
}

@Injectable({ providedIn: 'root' })
export class CustomerService {
  private apiUrl = 'http://localhost:8080/clientes';

  constructor(private http: HttpClient) {}

  getAll(): Observable<Customer[]> {
    return this.http.get<Customer[]>(this.apiUrl);
  }

  getById(id: string): Observable<Customer> {
    return this.http.get<Customer>(`${this.apiUrl}/${id}`);
  }

  create(payload: { name: string; email: string }): Observable<Customer> {
    return this.http.post<Customer>(this.apiUrl, payload);
  }

  update(id: string, payload: { name: string; email: string }): Observable<Customer> {
    return this.http.put<Customer>(`${this.apiUrl}/${id}`, payload);
  }

  delete(id: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }

  deposit(id: string, amount: number): Observable<Customer> {
    return this.http.post<Customer>(`${this.apiUrl}/${id}/depositar`, { amount });
  }

  withdraw(id: string, amount: number): Observable<Customer> {
    return this.http.post<Customer>(`${this.apiUrl}/${id}/sacar`, { amount });
  }

  getTransactions(id: string, page = 1, size = 10): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}/${id}/transacoes?page=${page}&size=${size}`);
  }
}
