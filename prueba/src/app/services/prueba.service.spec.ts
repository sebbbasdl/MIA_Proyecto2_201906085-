import { TestBed } from '@angular/core/testing';

import { EjemploService } from './prueba.service';

describe('EjemploService', () => {
  let service: EjemploService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(EjemploService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
