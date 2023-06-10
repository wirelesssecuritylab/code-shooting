import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RangeListComponent } from './range-list.component';

describe('RangeListComponent', () => {
  let component: RangeListComponent;
  let fixture: ComponentFixture<RangeListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ RangeListComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(RangeListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
