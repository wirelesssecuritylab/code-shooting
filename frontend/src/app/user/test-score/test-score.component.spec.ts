import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TestScoreComponent } from './test-score.component';

describe('TestScoreComponent', () => {
  let component: TestScoreComponent;
  let fixture: ComponentFixture<TestScoreComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TestScoreComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TestScoreComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
