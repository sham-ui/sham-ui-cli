import { createDI } from 'sham-ui';
import * as directives from 'sham-ui-directives';
import renderer from 'sham-ui-test-helpers';
//eslint-disable-next-line max-len
import LayoutNavigationThemeSwitch  from '../../../../../src/components/layout/navigation/theme-switch.sfc';

it( 'renders correctly', () => {
    const meta = renderer( LayoutNavigationThemeSwitch, {}, {
        directives
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'toggle theme', () => {
    const DI = createDI();
    const toggle = jest.fn();
    DI.bind( 'dark-theme', {
        toggle: toggle
    } );
    const meta = renderer( LayoutNavigationThemeSwitch, {}, {
        DI,
        directives
    } );
    meta.ctx.container.querySelector( '#checkbox' ).click();
    expect( meta.toJSON() ).toMatchSnapshot();
    expect( toggle ).toHaveBeenCalledTimes( 1 );
} );
