from typing import List

class Company:
    def __init__(self, bce_number: str, name: str, address: str, nace_codes: List[str]) -> None:
        self.bce_number = bce_number
        self.name = name
        self.address = name
        self.nace_codes = nace_codes
    
    def __str__(self) -> str:
        return f"Entreprise {self.name} (BCE: {self.bce_number})"

    def is_tech_company(self) -> bool:
        tech_codes = {"62010", "62020", "62030", "62090"}

        return any(code in tech_codes for code in self.nace_codes)