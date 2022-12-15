import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesArticlesCategoriesPage  from '../../../../../../src/components/routes/articles/categories/page.sfc';
import renderer from 'sham-ui-test-helpers';
import { storage } from '../../../../../../src/storages/session';

it( 'renders correctly', () => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    storage( DI ).sessionValidated = true;

    const getMock = jest.fn();
    DI.bind( 'store', {
        articleCategories: getMock.mockReturnValueOnce(
            Promise.resolve( { meta: {} } )
        )
    } );

    const meta = renderer( RoutesArticlesCategoriesPage, {}, {
        DI,
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'display errors', async() => {
    const DI = createDI();
    DI.bind( 'title', {
        change() {}
    } );
    storage( DI ).sessionValidated = true;
    DI.bind( 'store', {
        articleCategories: jest.fn().mockReturnValueOnce(
            Promise.reject( {} )
        )
    } );
    const meta = renderer( RoutesArticlesCategoriesPage, {}, {
        DI,
        directives
    } );
    await new Promise( resolve => setImmediate( resolve ) );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
