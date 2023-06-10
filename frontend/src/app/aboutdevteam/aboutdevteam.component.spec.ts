import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AboutdevteamComponent } from './aboutdevteam.component';

describe('AboutdevteamComponent', () => {
  let component: AboutdevteamComponent;
  let fixture: ComponentFixture<AboutdevteamComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ AboutdevteamComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(AboutdevteamComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
