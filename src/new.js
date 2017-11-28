const path = require( 'path' );
const Metalsmith = require( 'metalsmith' );
const getOptions = require( './options/new' );

const filterFiles = require( './metalsmith-plugins/filter-files' );
const promptQuestions = require( './metalsmith-plugins/prompt' );
const renderTemplates = require( './metalsmith-plugins/render-templates' );

function scaffold( name, src, dest, done ) {
    const opts = getOptions( name, src );
    const metalsmith = Metalsmith( path.join( src, 'template' ) );

    const data = Object.assign( metalsmith.metadata(), {
        destDirName: name,
        inPlace: dest === process.cwd(),
        noEscape: true
    } );

    metalsmith.use( promptQuestions( opts.prompts ) )
        .use( filterFiles( opts.filters ) )
        .use( renderTemplates( opts.skipInterpolation ) );

    metalsmith.clean( false )
        .source( '.' ) // start from template root instead of `./src` which is Metalsmith's default for `source`
        .destination( dest )
        .build( done );

    return data;
}


module.exports = scaffold;