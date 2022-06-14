import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
// eslint-disable-next-line max-len
import RoutesArticlesTagsPage  from '../../../../../../src/components/routes/articles/tags/page.sfc';
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
        articleTags: getMock.mockReturnValueOnce(
            Promise.resolve( { meta: {} } )
        )
    } );

    const meta = renderer( RoutesArticlesTagsPage, {
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
        articleTags: jest.fn().mockReturnValueOnce(
            Promise.reject( {} )
        )
    } );
    const meta = renderer( RoutesArticlesTagsPage, {
        DI,
        directives
    } );
    await new Promise( resolve => setImmediate( resolve ) );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
