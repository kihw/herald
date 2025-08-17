import React, { useState, useEffect } from "react";
import {
  Paper,
  TextField,
  Button,
  Typography,
  Box,
  Alert,
  CircularProgress,
  MenuItem,
  Divider,
} from "@mui/material";
import { SportsEsports, Person, Public } from "@mui/icons-material";
import { useAuth } from "../../context/AuthContext";
import { apiService, Region } from "../../services/api";

export function RiotValidationForm() {
  const { state, validateAccount, clearError } = useAuth();
  const [formData, setFormData] = useState({
    riotId: "",
    riotTag: "",
    region: "",
  });
  const [regions, setRegions] = useState<Region[]>([]);
  const [loadingRegions, setLoadingRegions] = useState(true);

  // Mapping des codes de région vers les noms affichables
  const getRegionName = (code: string): string => {
    const regionNames: Record<string, string> = {
      br1: "Brazil",
      eun1: "Europe Nordic & East",
      euw1: "Europe West",
      jp1: "Japan",
      kr: "Korea",
      la1: "Latin America North",
      la2: "Latin America South",
      na1: "North America",
      oc1: "Oceania",
      tr1: "Turkey",
      ru: "Russia",
    };
    return regionNames[code] || code.toUpperCase();
  };

  // Charger les régions supportées au démarrage
  useEffect(() => {
    const loadRegions = async () => {
      try {
        const response = await apiService.getSupportedRegions();
        // Les régions sont déjà formatées par le serveur réel
        const formattedRegions = response.regions.map((region: any) =>
          typeof region === "string"
            ? { code: region, name: getRegionName(region) }
            : { code: region.code, name: region.name }
        );
        setRegions(formattedRegions);
        // Définir EUW1 comme région par défaut si disponible
        if (formattedRegions.some((r) => r.code === "euw1")) {
          setFormData((prev) => ({ ...prev, region: "euw1" }));
        } else if (formattedRegions.length > 0) {
          setFormData((prev) => ({
            ...prev,
            region: formattedRegions[0].code,
          }));
        }
      } catch (error) {
        console.error("Erreur lors du chargement des régions:", error);
        // Utiliser des régions par défaut en cas d'erreur
        const defaultRegions = [
          { code: "euw1", name: "Europe West" },
          { code: "eun1", name: "Europe Nordic & East" },
          { code: "na1", name: "North America" },
          { code: "kr", name: "Korea" },
          { code: "br1", name: "Brazil" },
          { code: "la1", name: "Latin America North" },
          { code: "la2", name: "Latin America South" },
          { code: "oc1", name: "Oceania" },
          { code: "tr1", name: "Turkey" },
          { code: "ru", name: "Russia" },
          { code: "jp1", name: "Japan" },
        ];
        setRegions(defaultRegions);
        setFormData((prev) => ({ ...prev, region: "euw1" }));
      } finally {
        setLoadingRegions(false);
      }
    };

    loadRegions();
  }, []);

  const handleInputChange =
    (field: string) => (event: React.ChangeEvent<HTMLInputElement>) => {
      setFormData((prev) => ({
        ...prev,
        [field]: event.target.value,
      }));

      // Effacer l'erreur quand l'utilisateur commence à taper
      if (state.error) {
        clearError();
      }
    };

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    // Validation des champs obligatoires
    if (!formData.riotId.trim()) {
      console.error("Riot ID est requis");
      return;
    }
    if (!formData.riotTag.trim()) {
      console.error("Riot Tag est requis");
      return;
    }
    if (!formData.region.trim()) {
      console.error("Région est requise");
      return;
    }

    await validateAccount(
      formData.riotId.trim(),
      formData.riotTag.trim(),
      formData.region
    );
  };

  return (
    <Paper elevation={3} sx={{ p: 4 }}>
      <Box sx={{ textAlign: "center", mb: 3 }}>
        <SportsEsports sx={{ fontSize: 48, color: "primary.main", mb: 2 }} />
        <Typography variant="h4" component="h1" gutterBottom>
          LoL Match Manager
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Connectez-vous avec votre compte Riot Games
        </Typography>
      </Box>

      <Divider sx={{ my: 3 }} />

      <form onSubmit={handleSubmit}>
        <Box sx={{ display: "flex", flexDirection: "column", gap: 3 }}>
          {/* Riot ID */}
          <TextField
            fullWidth
            label="Riot ID"
            placeholder="VotreNomDeJoueur"
            value={formData.riotId}
            onChange={handleInputChange("riotId")}
            disabled={state.isLoading}
            required
            InputProps={{
              startAdornment: <Person sx={{ color: "action.active", mr: 1 }} />,
            }}
            helperText="Votre nom d'invocateur sans le tag (ex: Player123)"
          />

          {/* Riot Tag */}
          <TextField
            fullWidth
            label="Riot Tag"
            placeholder="EUW1"
            value={formData.riotTag}
            onChange={handleInputChange("riotTag")}
            disabled={state.isLoading}
            required
            helperText="Votre tag Riot sans le # (ex: EUW1, NA1, 1234)"
          />

          {/* Région */}
          <TextField
            fullWidth
            select
            label="Région"
            value={formData.region}
            onChange={handleInputChange("region")}
            disabled={state.isLoading || loadingRegions}
            required
            InputProps={{
              startAdornment: <Public sx={{ color: "action.active", mr: 1 }} />,
            }}
            helperText="Sélectionnez votre région de jeu"
          >
            {loadingRegions ? (
              <MenuItem value="">
                <CircularProgress size={20} sx={{ mr: 1 }} />
                Chargement...
              </MenuItem>
            ) : (
              regions.map((region) => (
                <MenuItem key={region.code} value={region.code}>
                  {region.name}
                </MenuItem>
              ))
            )}
          </TextField>

          {/* Message d'erreur */}
          {state.error && (
            <Alert severity="error" onClose={clearError}>
              {state.error}
            </Alert>
          )}

          {/* Bouton de validation */}
          <Button
            type="submit"
            fullWidth
            variant="contained"
            size="large"
            disabled={
              state.isLoading ||
              !formData.riotId.trim() ||
              !formData.riotTag.trim() ||
              loadingRegions
            }
            sx={{ mt: 2, py: 1.5 }}
          >
            {state.isLoading ? (
              <>
                <CircularProgress size={20} sx={{ mr: 1 }} />
                Validation en cours...
              </>
            ) : (
              "Valider le compte"
            )}
          </Button>
        </Box>
      </form>

      <Divider sx={{ my: 3 }} />

      <Box sx={{ textAlign: "center" }}>
        <Typography variant="body2" color="text.secondary">
          Exemple: Riot ID "Canna", Tag "KC", Région "Europe West"
        </Typography>
        <Typography
          variant="caption"
          color="text.secondary"
          sx={{ mt: 1, display: "block" }}
        >
          Nous validons votre compte via l'API officielle Riot Games
        </Typography>
      </Box>
    </Paper>
  );
}
