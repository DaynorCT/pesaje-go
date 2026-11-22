export namespace main {
	
	export class EstadoIntegradoDTO {
	    balanzaConectada: boolean;
	    balanzaTexto: string;
	    syncTexto: string;
	    syncVariante: string;
	    pesoVivo: number;
	    pesoCapturado: number;
	    textoConectar: string;
	    puedeCapturar: boolean;
	    puedeEnviar: boolean;
	    syncDisponible: boolean;
	
	    static createFrom(source: any = {}) {
	        return new EstadoIntegradoDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.balanzaConectada = source["balanzaConectada"];
	        this.balanzaTexto = source["balanzaTexto"];
	        this.syncTexto = source["syncTexto"];
	        this.syncVariante = source["syncVariante"];
	        this.pesoVivo = source["pesoVivo"];
	        this.pesoCapturado = source["pesoCapturado"];
	        this.textoConectar = source["textoConectar"];
	        this.puedeCapturar = source["puedeCapturar"];
	        this.puedeEnviar = source["puedeEnviar"];
	        this.syncDisponible = source["syncDisponible"];
	    }
	}
	export class AvisoDTO {
	    titulo: string;
	    mensaje: string;
	    esError: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AvisoDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.titulo = source["titulo"];
	        this.mensaje = source["mensaje"];
	        this.esError = source["esError"];
	    }
	}
	export class AccionDTO {
	    ok: boolean;
	    aviso?: AvisoDTO;
	    estado?: EstadoIntegradoDTO;
	
	    static createFrom(source: any = {}) {
	        return new AccionDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.aviso = this.convertValues(source["aviso"], AvisoDTO);
	        this.estado = this.convertValues(source["estado"], EstadoIntegradoDTO);
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

