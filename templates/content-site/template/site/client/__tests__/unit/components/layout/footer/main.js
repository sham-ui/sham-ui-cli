import LayoutFooterMain  from '../../../../../src/components/layout/footer/main.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( LayoutFooterMain, {} );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
