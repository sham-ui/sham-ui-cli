const { Compiler } = require( 'sham-ui-templates' );
const { parser } = require( 'sham-ui-templates/lib/parser' );
const { visit } = require( 'sham-ui-templates/lib/visitor' );
const { promisify } = require( 'util' );
const path = require( 'path' );
const fs = require( 'fs' );

const readdir = promisify( fs.readdir );
const stat = promisify( fs.stat );
const transforms = new Compiler().transforms;

const VISITORS = {
    filters( items ) {
        return {
            FilterExpression( node ) {
                items.add( node.callee.name );
            }
        }
    },
    directives( items ) {
        return {
            Directive( node ) {
                items.add( node.name );
            }
        }
    },
};

const OPTIONS_BY_EXTENSIONS = {
    '.sht': {},
    '.sfc': {
        asModule: false,
        asSingleFileComponent: true,
    }
};

async function getFiles( dir ) {
    const subdirs = await readdir( dir );
    const files = await Promise.all( subdirs.map( async ( subdir ) => {
        const res = path.resolve( dir, subdir );
        return ( await stat( res ) ).isDirectory() ?
            getFiles( res ) :
            res
        ;
    } ) );
    return files.reduce( ( a, x ) => a.concat( x ), [] );
}

function parse( filePath ) {
    const code = fs.readFileSync( filePath ).toString();
    const ast = parser.parse( filePath, code );
    const options = OPTIONS_BY_EXTENSIONS[ path.extname( filePath ) ];
    transforms.forEach( transform => transform( ast, options ) );
    return ast;

}

async function collect( type, projectFiles ) {
    if ( !( type in VISITORS ) ) {
        throw new Error( 'Unknown type: ' + type );
    }
    const extensions = Object.keys( OPTIONS_BY_EXTENSIONS );
    const files = ( await getFiles( projectFiles ) )
        .filter(
            x => extensions.includes( path.extname( x ) )
        );
    const foundItems = new Set();
    const visitor = VISITORS[ type ]( foundItems );
    files.forEach( file => {
        const ast = parse( file );
        visit( ast, visitor )
    } );
    return [ ...foundItems ];
}

module.exports = collect;
