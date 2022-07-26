import { babel } from '@rollup/plugin-babel';
import commonjs from '@rollup/plugin-commonjs';
import replace from '@rollup/plugin-replace';
import cleaner from 'rollup-plugin-cleaner';
import scss from 'rollup-plugin-scss';
import serve from 'rollup-plugin-dev';
import { terser } from 'rollup-plugin-terser';
import nodeResolveWithMacro from 'rollup-plugin-node-resolve-with-sham-ui-macro';
import shamUICompiler from 'rollup-plugin-sham-ui-templates';
import copy from 'rollup-plugin-copy';
import pkg from './package.json';

const prod = !process.env.ROLLUP_WATCH;
const dev = !prod;

const superUserPages = [
    'members',
    'server_info'
];

export default {
    input: 'src/main.js',
    preserveEntrySignatures: false,
    output: {
        dir: 'dist',
        format: 'system',
        sourcemap: true,
        entryFileNames: 'bundle.js',
        chunkFileNames( chunkInfo ) {
            const id = chunkInfo.facadeModuleId;
            if ( id && id.endsWith( 'page.sfc' ) ) {
                const page = id.replace( '/page.sfc', '' )
                    .split( '/' )
                    .pop()
                    .replace( /-/g, '_' );
                if ( superUserPages.includes( page ) ) {
                    return `su_${page}.js`;
                }
                return `${page}.js`;
            }
            return '[name].js';
        }
    },
    plugins: [
        cleaner( {
            targets: [
                'dist'
            ]
        } ),
        copy( {
            targets: [
                { src: 'index.html', dest: 'dist' },
                { src: 'favicon.ico', dest: 'dist' },
                { src: './node_modules/systemjs/dist/s.min.js', dest: 'dist' },
                {
                    src: [
                        'src/font/fontello.eot',
                        'src/font/fontello.svg',
                        'src/font/fontello.ttf',
                        'src/font/fontello.woff',
                        'src/font/fontello.woff2'
                    ],
                    dest: 'dist/font'
                }
            ]
        } ),
        replace( {
            preventAssignment: true,
            'process.env.NODE_ENV': JSON.stringify( process.env.NODE_ENV ),
            PRODUCTION: JSON.stringify( prod ),
            VERSION: JSON.stringify( pkg.version )
        } ),
        shamUICompiler( {
            extensions: [ '.sht' ],
            compilerOptions: {
                removeDataTest: prod
            }
        } ),
        shamUICompiler( {
            extensions: [ '.sfc' ],
            compilerOptions: {
                asModule: false,
                asSingleFileComponent: true,
                removeDataTest: prod
            }
        } ),
        nodeResolveWithMacro( {
            browser: true
        } ),
        babel( {
            extensions: [ '.js', '.sht', '.sfc' ],
            exclude: [ 'node_modules/**' ],
            babelHelpers: 'bundled'
        } ),
        commonjs(),
        scss( {
            failOnError: true,
            output: 'dist/bundle.css',
            outputStyle: 'compressed',
            watch: 'src/styles',
            sass: require( 'sass' )
        } ),
        dev && serve( {
            port: 3000,
            dirs: [ 'dist' ],
            spa: true,
            proxy: [
                { from: '/assets', to: 'http://localhost:3003/assets' }
            ]
        } ),
        prod && terser()
    ]
};
