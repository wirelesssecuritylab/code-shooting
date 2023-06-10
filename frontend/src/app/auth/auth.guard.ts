import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot, UrlTree } from '@angular/router';
import { Observable } from 'rxjs';
import { Router } from '@angular/router';

const USEROLE: string[] = ['shootuser', 'admin'];
const LOGIN_URL = '/login';

@Injectable({
  providedIn: 'root'
})
export class AuthGuard implements CanActivate {

  constructor(public router: Router) {
  }

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
      const userNumber = localStorage.getItem('user');
      const userRole = localStorage.getItem('role');

      if (userNumber && USEROLE.indexOf(userRole) >= 0) {
        return true;
      }

      this.router.navigateByUrl(LOGIN_URL);
      return false;
  }

}
