import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/lib/href-to';
import renderer from 'sham-ui-test-helpers';
// eslint-disable-next-line max-len
import RoutesArticlesListPage  from '../../../../../../src/components/routes/articles/list/page.sfc';
import { storage } from '../../../../../../src/storages/session';

it( 'renders correctly', async() => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    storage( DI ).sessionValidated = true;
    DI.bind( 'router', {
        generate: jest.fn().mockReturnValue( '/' )
    } );

    DI.bind( 'store', {
        articleCategories: jest.fn().mockReturnValueOnce(
            Promise.resolve( {
                categories: [
                    { 'id': '4', 'name': 'Hello world!', 'slug': 'hello-world' },
                    { 'id': '6', 'name': 'Название категории', 'slug': 'nazvanie-kategorii' }
                ],
                meta: {}
            } )
        ),
        articles: jest.fn().mockReturnValueOnce(
            Promise.resolve( {
                'articles': [
                    {
                        'id': '10',
                        'title': 'Тест2',
                        'slug': 'test2',
                        'category_id': '1',
                        'published_at': '2022-04-25T19:34:23.619+07:00'
                    }
                ],
                'meta': {
                    'limit': 50, 'offset': 0, 'total': 1
                }
            } )
        )
    } );

    const meta = renderer( RoutesArticlesListPage, {}, {
        DI,
        directives: {
            ...directives,
            hrefto
        }
    } );
    await new Promise( resolve => setImmediate( resolve ) );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'fail load categories data', async() => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    storage( DI ).sessionValidated = true;
    DI.bind( 'router', {
        generate: jest.fn().mockReturnValue( '/' )
    } );

    DI.bind( 'store', {
        articleCategories: jest.fn().mockReturnValueOnce(
            Promise.reject( {} )
        )
    } );

    const meta = renderer( RoutesArticlesListPage, {}, {
        DI,
        directives: {
            ...directives,
            hrefto
        }
    } );
    await new Promise( resolve => setImmediate( resolve ) );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'fail load articles', async() => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    storage( DI ).sessionValidated = true;
    DI.bind( 'router', {
        generate: jest.fn().mockReturnValue( '/' )
    } );

    DI.bind( 'store', {
        articleCategories: jest.fn().mockReturnValueOnce(
            Promise.resolve( {
                categories: [
                    { 'id': '4', 'name': 'Hello world!', 'slug': 'hello-world' },
                    { 'id': '6', 'name': 'Название категории', 'slug': 'nazvanie-kategorii' }
                ],
                meta: {}
            } )
        ),
        articles: jest.fn().mockReturnValueOnce(
            Promise.reject( {} )
        )
    } );

    const meta = renderer( RoutesArticlesListPage, {}, {
        DI,
        directives: {
            ...directives,
            hrefto
        }
    } );
    await new Promise( resolve => setImmediate( resolve ) );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
