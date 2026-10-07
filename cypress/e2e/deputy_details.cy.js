describe("Deputy details tab", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/lay/1");
    });

    it("has header for section", () => {
        cy.get("h1").should("contain", "Deputy details");
    });

    it("has a button which to allow the user to record deputy of death", () => {
        cy.get("#record-death").should("contain", "Record death of deputy").click();
    });

    it("lists active cases", () => {
        cy.get("#overview").should("contain", "3");
        cy.get("#overview").should("contain", "Active clients");
    });

    describe("Deputy contact details", () => {
        it("has rows in tables with accurate keys and values", () => {
            cy.get("#deputy-details > :nth-child(1) > .govuk-summary-list__key").should("contain", "Name");
            cy.get("#deputy-details > :nth-child(1) > .govuk-summary-list__value").should("contain", "Mr Mike Deputy");

            cy.get("#deputy-details > :nth-child(2) > .govuk-summary-list__key").should("contain", "Other names");
            cy.get("#deputy-details > :nth-child(2) > .govuk-summary-list__value").should("contain", "Michael");

            cy.get("#deputy-details > :nth-child(3) > .govuk-summary-list__key").should("contain", "Date of birth");
            cy.get("#deputy-details > :nth-child(3) > .govuk-summary-list__value").should("contain", "01/01/1980");

            cy.get("#deputy-details > :nth-child(4) > .govuk-summary-list__key").should("contain", "Mobile telephone number");
            cy.get("#deputy-details > :nth-child(4) > .govuk-summary-list__value").should("contain", "0771 2345678");

            cy.get("#deputy-details > :nth-child(5) > .govuk-summary-list__key").should("contain", "Daytime telephone number");
            cy.get("#deputy-details > :nth-child(5) > .govuk-summary-list__value").should("contain", "0115 876 5574");


            cy.get("#deputy-details > :nth-child(6) > .govuk-summary-list__key").should("contain", "Evening telephone");
            cy.get("#deputy-details > :nth-child(6) > .govuk-summary-list__value").should("contain", "0115 2767825");

            cy.get("#deputy-details > :nth-child(7) > .govuk-summary-list__key").should("contain", "Email address");
            cy.get("#deputy-details > :nth-child(7) > .govuk-summary-list__value").should(
                "contain",
                "testemail@hotmail.co.uk",
            );
            cy.get("#deputy-details > :nth-child(8) > .govuk-summary-list__key").should("contain", "Postal address");
            cy.get("#deputy-details > :nth-child(8) > .govuk-summary-list__value")
                .should("contain.text", "Seax House")
                .and("contain.text", "19 Market Rd")
                .and("contain.text", "Chelmsford")
                .and("contain.text", "Essex")
                .and("contain.text", "CM1 1GG");

            cy.get("#deputy-details > :nth-child(9) > .govuk-summary-list__key").should("contain", "Is airmail required?");
            cy.get("#deputy-details > :nth-child(9) > .govuk-summary-list__value").should("contain", "No");
        });
    });

    describe("Additional information", () => {
        it("has rows in tables with accurate keys and values", () => {
            cy.get(":nth-child(2) > .govuk-summary-list > :nth-child(1) > .govuk-summary-list__key").should(
                "contain",
                "Special correspondence requirements",
            );
            cy.get(":nth-child(2) > .govuk-summary-list > :nth-child(1) > .govuk-summary-list__value").should(
                "contain",
                "Audio, Large Print",
            );
            cy.get(":nth-child(2) > .govuk-summary-list > :nth-child(2) > .govuk-summary-list__key").should(
                "contain",
                "Interpreter requirements",
            );
            cy.get(":nth-child(2) > .govuk-summary-list > :nth-child(3) > .govuk-summary-list__key").should(
                "contain",
                "Additional information",
            );
        });
    });
});
