import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Reportes3Component } from './reportes3.component';

describe('Reportes3Component', () => {
  let component: Reportes3Component;
  let fixture: ComponentFixture<Reportes3Component>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ Reportes3Component ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Reportes3Component);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
