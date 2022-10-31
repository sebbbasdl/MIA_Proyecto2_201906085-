import { Component, OnInit } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';
import { EjemploService } from 'src/app/services/prueba.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-reportes',
  templateUrl: './reportes3.component.html',
  styleUrls: ['./reportes3.component.css']
})
export class Reportes3Component implements OnInit {
  entrada = "";
  salida = "";

  imagePath: any;

  constructor(private _sanitizer: DomSanitizer, public service: EjemploService,private router: Router) { }

  ngOnInit(): void {
    this.service.getReporte3().subscribe((res:any) => {
      let img = JSON.parse(JSON.stringify(res.result))
      this.imagePath = this._sanitizer.bypassSecurityTrustResourceUrl(img);
    });
  }
  logout(){
    const cmd = "logout";
    console.log("---cmd----")
    console.log(cmd)

      this.service.postEntrada(cmd).subscribe(async (res:any) => {
        this.salida += await res.result + "\n";
        console.log("salida: "+this.salida)
        this.router.navigate(['login']);

      });


  }

}
