import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ActionTabComponent } from './action-tab.component';

describe('ActionTabComponent', () => {
  let component: ActionTabComponent;
  let fixture: ComponentFixture<ActionTabComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ActionTabComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ActionTabComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
