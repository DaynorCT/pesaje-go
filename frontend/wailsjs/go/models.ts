export namespace main {

	export class AbmFilaDTO {
	    id!: number;
	    estado!: string;
	    celdas!: Record<string, string>;
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        Object.assign(this, source);
	    }
	}
	export class AbmGuardarDTO {
	    entidad!: string;
	    id!: number;
	    campos!: Record<string, string>;
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        Object.assign(this, source);
	    }
	}
	export class AbmResultadoDTO {
	    ok!: boolean;
	    mensaje!: string;
	    esError!: boolean;
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        Object.assign(this, source);
	    }
	}
	export class ResumenInicioDTO {
	    pesajesTotal!: number;
	    pesajesHoy!: number;
	    vehiculosActivos!: number;
	    productosActivos!: number;
	    choferesActivos!: number;
	    proveedoresActivos!: number;
	    syncConfigurada!: boolean;
	    syncEstado!: string;
	    syncPendientes!: number;
	    syncErrores!: number;
	    constructor(source: any = {}) {
	        Object.assign(this, source);
	    }
	}
	export class BoletaConfigDTO {
	    id!: number;
	    salida!: string;
	    direccion!: string;
	    celular1!: string;
	    celular2!: string;
	    ciudad!: string;
	    tituloSuperior!: string;
	    marca!: string;
	    tituloDocumento!: string;
	    acreditacion!: string;
	    constructor(source: any = {}) {
	        Object.assign(this, source);
	    }
	}
	export class FiltroReporteDTO {
	    desde!: string;
	    hasta!: string;
	    placa!: string;
	    constructor(source: any = {}) {
	        Object.assign(this, source);
	    }
	}
	export class PdfResultadoDTO {
	    ok!: boolean;
	    ruta!: string;
	    mensaje!: string;
	    esError!: boolean;
	    constructor(source: any = {}) {
	        Object.assign(this, source);
	    }
	}
	export class BalanzaConfigDTO {
	    puertoCom!: string;
	    baudRate!: number;
	    dataBits!: number;
	    paridad!: string;
	    stopBits!: number;
	    flowControl!: string;
	    timeoutMs!: number;
	    metodoLectura!: string;
	    posicionInicio!: number;
	    posicionLongitud!: number;
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        Object.assign(this, source);
	    }
	}
	
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
	export class EstadoPesajePaneDTO {
	    balanzaConectada!: boolean;
	    textoConectar!: string;
	    pesoVivo!: number;
	    pesoBruto!: number;
	    pesoTara!: number;
	    pesoNeto!: number;
	    puedeCapturar!: boolean;
	    mensaje!: string;
	    mensajeEsError!: boolean;

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        Object.assign(this, source);
	    }
	}
	export class AccionDTO {
	    ok: boolean;
	    aviso?: AvisoDTO;
	    estado?: EstadoIntegradoDTO;
	    estadoPesaje?: EstadoPesajePaneDTO;
	
	    static createFrom(source: any = {}) {
	        return new AccionDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.aviso = this.convertValues(source["aviso"], AvisoDTO);
	        this.estado = this.convertValues(source["estado"], EstadoIntegradoDTO);
	        this.estadoPesaje = this.convertValues(source["estadoPesaje"], EstadoPesajePaneDTO);
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
	    personaId: number;
	    email: string;
	    nombre: string;
	    rol: string;
	
	    static createFrom(source: any = {}) {
	        return new UsuarioDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.personaId = source["personaId"];
	        this.email = source["email"];
	        this.nombre = source["nombre"];
	        this.rol = source["rol"];
	    }
	}
	export class OpcionDTO {
	    id: number;
	    label: string;
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source['id'];
	        this.label = source['label'];
	    }
	}
	export class PesajeOpcionesDTO {
	    vehiculos: OpcionDTO[];
	    choferes: OpcionDTO[];
	    proveedores: OpcionDTO[];
	    productos: OpcionDTO[];
	    origenes: OpcionDTO[];
	    destinos: OpcionDTO[];
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vehiculos = source['vehiculos'] || [];
	        this.choferes = source['choferes'] || [];
	        this.proveedores = source['proveedores'] || [];
	        this.productos = source['productos'] || [];
	        this.origenes = source['origenes'] || [];
	        this.destinos = source['destinos'] || [];
	    }
	}
	export class VehiculoDetalleDTO {
	    tipo: string;
	    tara: number;
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tipo = source['tipo'];
	        this.tara = source['tara'];
	    }
	}
	export class ChoferDetalleDTO {
	    licencia: string;
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.licencia = source['licencia'];
	    }
	}
	export class FilaPesajeDTO {
	    id!: number;
	    guia!: string;
	    nombreChofer!: string;
	    placa!: string;
	    licencia!: string;
	    tipoVehiculo!: string;
	    nombreProveedor!: string;
	    producto!: string;
	    origen!: string;
	    destino!: string;
	    costoBs!: number;
	    fleteBs!: number;
	    fecha!: string;
	    hora!: string;
	    pesoBruto!: number;
	    tara!: number;
	    pesoNeto!: number;
	    estado!: string;
	    estadoBalanza!: string;
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        Object.assign(this, source);
	    }
	}
	export class RegistrarPesajeDTO {
	    vehiculoId!: number;
	    choferId?: number;
	    proveedorId?: number;
	    productoId?: number;
	    origenId?: number;
	    destinoId?: number;
	    pesoBruto!: number;
	    tara!: number;
	    costoBs!: number;
	    fleteBsPendiente!: number;
	    taraManual!: boolean;
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        Object.assign(this, source);
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

