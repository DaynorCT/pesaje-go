export namespace main {
	
	export class AbmFilaDTO {
	    id: number;
	    estado: string;
	    celdas: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new AbmFilaDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.estado = source["estado"];
	        this.celdas = source["celdas"];
	    }
	}
	export class AbmGuardarDTO {
	    entidad: string;
	    id: number;
	    campos: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new AbmGuardarDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entidad = source["entidad"];
	        this.id = source["id"];
	        this.campos = source["campos"];
	    }
	}
	export class AbmResultadoDTO {
	    ok: boolean;
	    mensaje: string;
	    esError: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AbmResultadoDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.mensaje = source["mensaje"];
	        this.esError = source["esError"];
	    }
	}
	export class EstadoPesajePaneDTO {
	    balanzaConectada: boolean;
	    textoConectar: string;
	    pesoVivo: number;
	    pesoBruto: number;
	    pesoTara: number;
	    pesoNeto: number;
	    puedeCapturar: boolean;
	    mensaje: string;
	    mensajeEsError: boolean;
	
	    static createFrom(source: any = {}) {
	        return new EstadoPesajePaneDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.balanzaConectada = source["balanzaConectada"];
	        this.textoConectar = source["textoConectar"];
	        this.pesoVivo = source["pesoVivo"];
	        this.pesoBruto = source["pesoBruto"];
	        this.pesoTara = source["pesoTara"];
	        this.pesoNeto = source["pesoNeto"];
	        this.puedeCapturar = source["puedeCapturar"];
	        this.mensaje = source["mensaje"];
	        this.mensajeEsError = source["mensajeEsError"];
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
	export class ActualizarPesajeDTO {
	    id: number;
	    vehiculoId: number;
	    choferId?: number;
	    proveedorId?: number;
	    productoId?: number;
	    origenId?: number;
	    destinoId?: number;
	    pesoBruto: number;
	    tara: number;
	    costoBs: number;
	    fleteBsPendiente: number;
	    taraManual: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ActualizarPesajeDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.vehiculoId = source["vehiculoId"];
	        this.choferId = source["choferId"];
	        this.proveedorId = source["proveedorId"];
	        this.productoId = source["productoId"];
	        this.origenId = source["origenId"];
	        this.destinoId = source["destinoId"];
	        this.pesoBruto = source["pesoBruto"];
	        this.tara = source["tara"];
	        this.costoBs = source["costoBs"];
	        this.fleteBsPendiente = source["fleteBsPendiente"];
	        this.taraManual = source["taraManual"];
	    }
	}
	
	export class BalanzaConfigDTO {
	    puertoCom: string;
	    baudRate: number;
	    dataBits: number;
	    paridad: string;
	    stopBits: number;
	    flowControl: string;
	    timeoutMs: number;
	    metodoLectura: string;
	    posicionInicio: number;
	    posicionLongitud: number;
	
	    static createFrom(source: any = {}) {
	        return new BalanzaConfigDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.puertoCom = source["puertoCom"];
	        this.baudRate = source["baudRate"];
	        this.dataBits = source["dataBits"];
	        this.paridad = source["paridad"];
	        this.stopBits = source["stopBits"];
	        this.flowControl = source["flowControl"];
	        this.timeoutMs = source["timeoutMs"];
	        this.metodoLectura = source["metodoLectura"];
	        this.posicionInicio = source["posicionInicio"];
	        this.posicionLongitud = source["posicionLongitud"];
	    }
	}
	export class BoletaConfigDTO {
	    id: number;
	    salida: string;
	    direccion: string;
	    celular1: string;
	    celular2: string;
	    ciudad: string;
	    tituloSuperior: string;
	    marca: string;
	    tituloDocumento: string;
	    acreditacion: string;
	
	    static createFrom(source: any = {}) {
	        return new BoletaConfigDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.salida = source["salida"];
	        this.direccion = source["direccion"];
	        this.celular1 = source["celular1"];
	        this.celular2 = source["celular2"];
	        this.ciudad = source["ciudad"];
	        this.tituloSuperior = source["tituloSuperior"];
	        this.marca = source["marca"];
	        this.tituloDocumento = source["tituloDocumento"];
	        this.acreditacion = source["acreditacion"];
	    }
	}
	export class ChoferDetalleDTO {
	    licencia: string;
	
	    static createFrom(source: any = {}) {
	        return new ChoferDetalleDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.licencia = source["licencia"];
	    }
	}
	
	
	export class FilaPesajeDTO {
	    id: number;
	    guia: string;
	    nombreChofer: string;
	    placa: string;
	    licencia: string;
	    tipoVehiculo: string;
	    nombreProveedor: string;
	    producto: string;
	    origen: string;
	    destino: string;
	    costoBs: number;
	    fleteBs: number;
	    fecha: string;
	    hora: string;
	    pesoBruto: number;
	    tara: number;
	    pesoNeto: number;
	    estado: string;
	    estadoBalanza: string;
	
	    static createFrom(source: any = {}) {
	        return new FilaPesajeDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.guia = source["guia"];
	        this.nombreChofer = source["nombreChofer"];
	        this.placa = source["placa"];
	        this.licencia = source["licencia"];
	        this.tipoVehiculo = source["tipoVehiculo"];
	        this.nombreProveedor = source["nombreProveedor"];
	        this.producto = source["producto"];
	        this.origen = source["origen"];
	        this.destino = source["destino"];
	        this.costoBs = source["costoBs"];
	        this.fleteBs = source["fleteBs"];
	        this.fecha = source["fecha"];
	        this.hora = source["hora"];
	        this.pesoBruto = source["pesoBruto"];
	        this.tara = source["tara"];
	        this.pesoNeto = source["pesoNeto"];
	        this.estado = source["estado"];
	        this.estadoBalanza = source["estadoBalanza"];
	    }
	}
	export class FiltroReporteDTO {
	    desde: string;
	    hasta: string;
	    placa: string;
	
	    static createFrom(source: any = {}) {
	        return new FiltroReporteDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.desde = source["desde"];
	        this.hasta = source["hasta"];
	        this.placa = source["placa"];
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
	export class OpcionDTO {
	    id: number;
	    label: string;
	
	    static createFrom(source: any = {}) {
	        return new OpcionDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	    }
	}
	export class PdfResultadoDTO {
	    ok: boolean;
	    ruta: string;
	    mensaje: string;
	    esError: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PdfResultadoDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.ruta = source["ruta"];
	        this.mensaje = source["mensaje"];
	        this.esError = source["esError"];
	    }
	}
	export class PesajeEdicionDTO {
	    id: number;
	    vehiculoId: number;
	    choferId?: number;
	    proveedorId?: number;
	    productoId?: number;
	    origenId?: number;
	    destinoId?: number;
	    pesoBruto: number;
	    tara: number;
	    costoBs: number;
	    fleteBsPendiente: number;
	    estadoBalanza: string;
	
	    static createFrom(source: any = {}) {
	        return new PesajeEdicionDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.vehiculoId = source["vehiculoId"];
	        this.choferId = source["choferId"];
	        this.proveedorId = source["proveedorId"];
	        this.productoId = source["productoId"];
	        this.origenId = source["origenId"];
	        this.destinoId = source["destinoId"];
	        this.pesoBruto = source["pesoBruto"];
	        this.tara = source["tara"];
	        this.costoBs = source["costoBs"];
	        this.fleteBsPendiente = source["fleteBsPendiente"];
	        this.estadoBalanza = source["estadoBalanza"];
	    }
	}
	export class PesajeOpcionesDTO {
	    vehiculos: OpcionDTO[];
	    choferes: OpcionDTO[];
	    proveedores: OpcionDTO[];
	    productos: OpcionDTO[];
	    origenes: OpcionDTO[];
	    destinos: OpcionDTO[];
	
	    static createFrom(source: any = {}) {
	        return new PesajeOpcionesDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vehiculos = this.convertValues(source["vehiculos"], OpcionDTO);
	        this.choferes = this.convertValues(source["choferes"], OpcionDTO);
	        this.proveedores = this.convertValues(source["proveedores"], OpcionDTO);
	        this.productos = this.convertValues(source["productos"], OpcionDTO);
	        this.origenes = this.convertValues(source["origenes"], OpcionDTO);
	        this.destinos = this.convertValues(source["destinos"], OpcionDTO);
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
	export class RegistrarPesajeDTO {
	    vehiculoId: number;
	    choferId?: number;
	    proveedorId?: number;
	    productoId?: number;
	    origenId?: number;
	    destinoId?: number;
	    pesoBruto: number;
	    tara: number;
	    costoBs: number;
	    fleteBsPendiente: number;
	    taraManual: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RegistrarPesajeDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vehiculoId = source["vehiculoId"];
	        this.choferId = source["choferId"];
	        this.proveedorId = source["proveedorId"];
	        this.productoId = source["productoId"];
	        this.origenId = source["origenId"];
	        this.destinoId = source["destinoId"];
	        this.pesoBruto = source["pesoBruto"];
	        this.tara = source["tara"];
	        this.costoBs = source["costoBs"];
	        this.fleteBsPendiente = source["fleteBsPendiente"];
	        this.taraManual = source["taraManual"];
	    }
	}
	export class ResumenInicioDTO {
	    pesajesTotal: number;
	    pesajesHoy: number;
	    vehiculosActivos: number;
	    productosActivos: number;
	    choferesActivos: number;
	    proveedoresActivos: number;
	    syncConfigurada: boolean;
	    syncEstado: string;
	    syncPendientes: number;
	    syncErrores: number;
	
	    static createFrom(source: any = {}) {
	        return new ResumenInicioDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pesajesTotal = source["pesajesTotal"];
	        this.pesajesHoy = source["pesajesHoy"];
	        this.vehiculosActivos = source["vehiculosActivos"];
	        this.productosActivos = source["productosActivos"];
	        this.choferesActivos = source["choferesActivos"];
	        this.proveedoresActivos = source["proveedoresActivos"];
	        this.syncConfigurada = source["syncConfigurada"];
	        this.syncEstado = source["syncEstado"];
	        this.syncPendientes = source["syncPendientes"];
	        this.syncErrores = source["syncErrores"];
	    }
	}
	
	export class VehiculoDetalleDTO {
	    tipo: string;
	    tara: number;
	
	    static createFrom(source: any = {}) {
	        return new VehiculoDetalleDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tipo = source["tipo"];
	        this.tara = source["tara"];
	    }
	}

}

