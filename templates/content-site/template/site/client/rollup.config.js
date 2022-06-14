import { babel } from '@rollup/plugin-babel';
import commonjs from '@rollup/plugin-commonjs';
import replace from '@rollup/plugin-replace';
import cleaner from 'rollup-plugin-cleaner';
import scss from 'rollup-plugin-scss';
import serve from 'rollup-plugin-serve';
import { terser } from 'rollup-plugin-terser';
import nodeResolveWithMacro from 'rollup-plugin-node-resolve-with-sham-ui-macro';
import shamUICompiler from 'rollup-plugin-sham-ui-templates';
import copy from 'rollup-plugin-copy';
import json from '@rollup/plugin-json';
import pkg from './package.json';

const prod = !process.env.ROLLUP_WATCH;

let projects = [];

if ( prod ) {
    projects = [ {
        input: 'src/ssr.js',
        output: {
            file: 'ssr/bundle.js',
            format: 'cjs',
            inlineDynamicImports: true
        },
        plugins: [
            cleaner( {
                targets: [
                    'ssr'
                ]
            } ),
            replace( {
                preventAssignment: true,
                'process.env.NODE_ENV': JSON.stringify( process.env.NODE_ENV ),
                PRODUCTION: JSON.stringify( prod ),
                VERSION: JSON.stringify( pkg.version ),
                IS_SSR: true
            } ),
            shamUICompiler( {
                extensions: [ '.sht' ]
            } ),
            shamUICompiler( {
                extensions: [ '.sfc' ],
                compilerOptions: {
                    asModule: false,
                    asSingleFileComponent: true
                }
            } ),
            nodeResolveWithMacro( {
                browser: false,
                preferBuiltins: true
            } ),
            babel( {
                extensions: [ '.js', '.sht', '.sfc' ],
                exclude: [ 'node_modules/**' ],
                babelHelpers: 'bundled'
            } ),
            commonjs(),
            json()
        ]
    }, {
        input: 'src/browser.js',
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
                    { src: 'favicon.ico', dest: 'dist' },
                    { src: 'src/images', dest: 'dist' },
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
                VERSION: JSON.stringify( pkg.version ),
                IS_SSR: false
            } ),
            shamUICompiler( {
                extensions: [ '.sht' ]
            } ),
            shamUICompiler( {
                extensions: [ '.sfc' ],
                compilerOptions: {
                    asModule: false,
                    asSingleFileComponent: true
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
            terser()
        ]
    } ];
} else {
    projects = [ {
        input: 'src/browser-dev.js',
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
                    { src: 'src/images', dest: 'dist' },
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
                VERSION: JSON.stringify( pkg.version ),
                IS_SSR: false
            } ),
            shamUICompiler( {
                extensions: [ '.sht' ],
                compilerOptions: {
                    removeDataTest: false
                }
            } ),
            shamUICompiler( {
                extensions: [ '.sfc' ],
                compilerOptions: {
                    asModule: false,
                    asSingleFileComponent: true,
                    removeDataTest: false
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
            serve( {
                port: 3000,
                contentBase: 'dist',
                historyApiFallback: '/index.html'
            } )
        ]
    } ];
}

export default projects;
