import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GradeBoardComponent } from './grade-board.component';

describe('GradeBoardComponent', () => {
  let component: GradeBoardComponent;
  let fixture: ComponentFixture<GradeBoardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ GradeBoardComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(GradeBoardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
