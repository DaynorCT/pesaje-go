export namespace main {
	
	export class UsuarioDTO {
	    id: number;
	    email: string;
	    nombre: string;
	    rol: string;
	
	    static createFrom(source: any = {}) {
	        return new UsuarioDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.email = source["email"];
	        this.nombre = source["nombre"];
	        this.rol = source["rol"];
	    }
	}
	export class LoginRespuesta {
	    exito: boolean;
	    estado: string;
	    mensaje: string;
	    usuario?: UsuarioDTO;
	
	    static createFrom(source: any = {}) {
	        return new LoginRespuesta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.exito = source["exito"];
	        this.estado = source["estado"];
	        this.mensaje = source["mensaje"];
	        this.usuario = this.convertValues(source["usuario"], UsuarioDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

