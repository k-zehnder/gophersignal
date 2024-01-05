describe('Articles Page', () => {
  it('should display the title "Latest Articles"', () => {
    cy.visit('http://localhost:3000');

    // Check for the presence of an h2 element with text "Latest Articles"
    cy.contains('h2', 'Latest Articles').should('be.visible');
  });
});
