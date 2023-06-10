import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProgressBoardComponent } from './progress-board.component';

describe('ProgressBoardComponent', () => {
  let component: ProgressBoardComponent;
  let fixture: ComponentFixture<ProgressBoardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ProgressBoardComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ProgressBoardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
