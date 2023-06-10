import { ComponentFixture, TestBed } from '@angular/core/testing';

import { QueryScoreComponent } from './query-score.component';

describe('QueryScoreComponent', () => {
  let component: QueryScoreComponent;
  let fixture: ComponentFixture<QueryScoreComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ QueryScoreComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(QueryScoreComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
