import { Component, OnInit } from '@angular/core';
import { RouteConfigLoadEnd } from '@angular/router';
import { EjemploService } from 'src/app/services/prueba.service';
import { Router } from '@angular/router';


@Component({
  selector: 'app-inicio',
  templateUrl: './inicio.component.html',
  styleUrls: ['./inicio.component.css']
})
export class InicioComponent implements OnInit {

  entrada = "";
  salida = "";


  constructor(public service: EjemploService) { }

  ngOnInit(): void {}

  public async onFileSelected(event:any) {
    const file:File = event.target.files[0];
    this.entrada = await file.text();
    console.log(this.entrada)
  }

  ejecutar(){
    this.salida = "--- Resultados ---\n";
    let split_entrada = this.entrada.split("\n");
    console.log(split_entrada)
    for (let i = 0; i < split_entrada.length; i++) {
      const cmd = split_entrada[i];
      console.log("---cmd----")
      console.log(cmd)
      if(cmd != ""){
        this.service.postEntrada(cmd).subscribe(async (res:any) => {
          this.salida += await res.result + "\n";
          console.log("salida:"+this.salida)

          //this.route.navigate(['/inicio'])
        });
      }

    }






  }

 /* ejecutar2(){
    const cmd = "login -usuario="+this.usuario+" -password="+this.pass+" -id="+this.id;
    console.log("---cmd----")
    console.log(cmd)
    if(cmd != ""){
      this.service.postEntrada(cmd).subscribe(async (res:any) => {
        this.salida += await res.result + "\n";
        console.log("salida: "+this.salida)
      });
    }
  }*/


}
