import { createDI } from 'sham-ui';
import * as diretives from 'sham-ui-directives';
import hrefto from 'sham-ui-router/lib/href-to';
import LayoutMain  from '../../../../src/components/layout/main.sfc';
import renderer, { compile } from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const DI = createDI();

    DI.bind( 'router', {
        generate: jest.fn().mockReturnValue( '/' ),
        activePageComponent: compile``,
        storage: {
            url: '/',
            addWatcher() {}
        }
    } );

    const meta = renderer( LayoutMain, {}, {
        DI,
        directives: {
            hrefto,
            ...diretives
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
