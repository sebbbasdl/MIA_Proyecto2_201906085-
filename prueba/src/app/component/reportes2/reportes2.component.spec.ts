import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Reportes2Component } from './reportes2.component';

describe('Reportes2Component', () => {
  let component: Reportes2Component;
  let fixture: ComponentFixture<Reportes2Component>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ Reportes2Component ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Reportes2Component);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
