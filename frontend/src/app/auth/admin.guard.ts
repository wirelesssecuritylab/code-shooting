import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot, UrlTree } from '@angular/router';
import { Observable } from 'rxjs';
import { Router } from '@angular/router';

const ADMINROLE: string = 'admin';
const LIST_URL: string = '/login'

@Injectable({
  providedIn: 'root'
})
export class AdminGuard implements CanActivate {

  constructor(private router: Router) {
  }

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
      const userNumber = localStorage.getItem('user');
      const userRole = localStorage.getItem('role');

      if (userNumber && userRole === ADMINROLE) {
        return true;
      }

      this.router.navigateByUrl(LIST_URL);
      return false;
  }
}
