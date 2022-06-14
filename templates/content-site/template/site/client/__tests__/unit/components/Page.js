import Page  from '../../../src/components/Page.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( Page );
    expect( meta.toJSON() ).toMatchSnapshot();
} );
