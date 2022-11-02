import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class EjemploService {

  constructor(
    private httpClient: HttpClient
  ) { }

  getNombreRandom(){
    return this.httpClient.get("https://randomuser.me/api/");
  }

  postEntrada(entrada: string){
    console.log("----------------------")
    console.log(entrada)
    console.log({ Cmd: entrada})
    console.log("**********************")


    return this.httpClient.post("http://52.14.179.229:5000/analizar",{ Cmd: entrada});
  }

  getReporte(){
    console.log("REPORTE")
    return this.httpClient.get("http://52.14.179.229:5000/reportes");
  }
  getReporte2(){
    return this.httpClient.get("http://52.14.179.229:5000/reportes2");
  }
  getReporte3(){
    return this.httpClient.get("http://52.14.179.229:5000/reportes3");
  }
}
