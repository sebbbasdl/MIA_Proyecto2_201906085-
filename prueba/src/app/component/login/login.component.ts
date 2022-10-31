import { Component, OnInit } from '@angular/core';
import { EjemploService } from 'src/app/services/prueba.service';
import { Router } from '@angular/router';
import { NavigationExtras } from '@angular/router';
@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  entrada = "";
  salida = "";
  pass="";
  id="";
  usuario="";


  constructor(public service: EjemploService,private router: Router) { }

  ngOnInit(): void {}

  public async onFileSelected(event:any) {
    const file:File = event.target.files[0];
    this.entrada = await file.text();
    console.log(this.entrada)
  }


  ejecutar2(){
    const cmd = "login -usuario="+this.usuario+" -password="+this.pass+" -id="+this.id;
    console.log("---cmd----")
    console.log(cmd)
    if(cmd != ""){
      this.service.postEntrada(cmd).subscribe(async (res:any) => {
        this.salida += await res.result + "\n";
        console.log("salida: "+this.salida)
        console.log(this.salida[0])
        if(this.salida[0]=="/"){
          console.log("entre")
          this.router.navigate(['reporte']);
        }else{
          this.router.navigate(['/login'])

        }
      });
    }

  }

  goToRoute() {

  }


}
