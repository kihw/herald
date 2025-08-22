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
  useTheme,
  Card,
  CardContent,
  Grid,
} from "@mui/material";
import { SportsEsports, Person, Public, Google } from "@mui/icons-material";
import { useAuth } from "../../context/AuthContext";
import { apiService, Region } from "../../services/api";
import { GoogleAuth } from "./GoogleAuth";
import { leagueColors } from "../../theme/leagueTheme";

export function RiotValidationForm() {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === 'dark';
  const { user, isAuthenticated, isLoading, error, validateAccount, clearError, securityStatus } = useAuth();
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
        const formattedRegions = (response.regions || [])
          .filter((region: any) => region) // Filtrer les éléments null/undefined
          .map((region: any) =>
            typeof region === "string"
              ? { code: region, name: getRegionName(region) }
              : { code: region?.code || '', name: region?.name || region?.code || 'Unknown' }
          )
          .filter((region: any) => region.code); // Filtrer les régions sans code
        setRegions(formattedRegions);
        // Définir EUW1 comme région par défaut si disponible
        if (formattedRegions.some((r) => r?.code === "euw1")) {
          setFormData((prev) => ({ ...prev, region: "euw1" }));
        } else if (formattedRegions.length > 0 && formattedRegions[0]?.code) {
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
      if (error) {
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

    validateAccount(
      formData.riotId.trim(),
      formData.riotTag.trim(),
      formData.region
    );
  };

  return (
    <Box 
      sx={{ 
        minHeight: '100vh',
        background: isDarkMode
          ? `linear-gradient(135deg, ${leagueColors.dark[50]} 0%, ${leagueColors.dark[100]} 100%)`
          : `linear-gradient(135deg, ${leagueColors.blue[50]} 0%, ${leagueColors.gold[50]} 100%)`,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        p: 2,
      }}
    >
      <Box sx={{ maxWidth: 500, width: '100%' }}>
        {/* Header */}
        <Card 
          elevation={0}
          sx={{ 
            mb: 3,
            background: isDarkMode
              ? `linear-gradient(135deg, ${leagueColors.blue[900]} 0%, ${leagueColors.dark[100]} 100%)`
              : `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
            color: '#fff',
            borderRadius: 3,
            textAlign: 'center',
          }}
        >
          <CardContent sx={{ py: 4 }}>
            <SportsEsports sx={{ fontSize: 64, color: leagueColors.gold[400], mb: 2 }} />
            <Typography variant="h3" component="h1" gutterBottom sx={{ fontWeight: 700, letterSpacing: 1 }}>
              Herald.lol
            </Typography>
            <Typography variant="h6" sx={{ opacity: 0.9, fontWeight: 400 }}>
              Analysez vos performances League of Legends
            </Typography>
          </CardContent>
        </Card>

        <Grid container spacing={3}>
          {/* Riot Games Authentication */}
          <Grid item xs={12}>
            <Card 
              elevation={0}
              sx={{ 
                borderRadius: 3,
                border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
                background: isDarkMode
                  ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
                  : `linear-gradient(135deg, #ffffff 0%, ${leagueColors.blue[25]} 100%)`,
              }}
            >
              <CardContent sx={{ p: 4 }}>
                <Typography 
                  variant="h5" 
                  gutterBottom 
                  sx={{ 
                    fontWeight: 600,
                    textAlign: 'center',
                    mb: 3,
                    color: 'primary.main',
                  }}
                >
                  Connexion Riot Games
                </Typography>

                <form onSubmit={handleSubmit}>
                  <Box sx={{ display: "flex", flexDirection: "column", gap: 3 }}>
                    {/* Riot ID */}
                    <TextField
                      fullWidth
                      label="Riot ID"
                      placeholder="VotreNomDeJoueur"
                      value={formData.riotId}
                      onChange={handleInputChange("riotId")}
                      disabled={isLoading}
                      required
                      InputProps={{
                        startAdornment: <Person sx={{ color: "primary.main", mr: 1 }} />,
                      }}
                      helperText="Votre nom d'invocateur sans le tag (ex: Player123)"
                      sx={{
                        '& .MuiOutlinedInput-root': {
                          borderRadius: 2,
                          '&:hover fieldset': {
                            borderColor: 'primary.main',
                          },
                        },
                      }}
                    />

                    {/* Riot Tag */}
                    <TextField
                      fullWidth
                      label="Riot Tag"
                      placeholder="EUW1"
                      value={formData.riotTag}
                      onChange={handleInputChange("riotTag")}
                      disabled={isLoading}
                      required
                      helperText="Votre tag Riot sans le # (ex: EUW1, NA1, 1234)"
                      sx={{
                        '& .MuiOutlinedInput-root': {
                          borderRadius: 2,
                          '&:hover fieldset': {
                            borderColor: 'primary.main',
                          },
                        },
                      }}
                    />

                    {/* Région */}
                    <TextField
                      fullWidth
                      select
                      label="Région"
                      value={formData.region}
                      onChange={handleInputChange("region")}
                      disabled={isLoading || loadingRegions}
                      required
                      InputProps={{
                        startAdornment: <Public sx={{ color: "primary.main", mr: 1 }} />,
                      }}
                      helperText="Sélectionnez votre région de jeu"
                      sx={{
                        '& .MuiOutlinedInput-root': {
                          borderRadius: 2,
                          '&:hover fieldset': {
                            borderColor: 'primary.main',
                          },
                        },
                      }}
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

                    {/* Message d'erreur avec surveillance sécurisée */}
                    {error && (
                      <Alert 
                        severity={securityStatus.suspiciousActivity ? "warning" : "error"}
                        onClose={clearError}
                        sx={{ 
                          borderRadius: 2,
                          border: `1px solid ${securityStatus.suspiciousActivity ? leagueColors.gold[400] : leagueColors.loss}`,
                        }}
                      >
                        {securityStatus.suspiciousActivity && (
                          <>
                            🚨 Activité suspecte détectée ({securityStatus.failedAttempts} tentatives) - 
                          </>
                        )}
                        {error}
                      </Alert>
                    )}

                    {/* Bouton de validation */}
                    <Button
                      type="submit"
                      fullWidth
                      variant="contained"
                      size="large"
                      disabled={
                        isLoading ||
                        !formData.riotId.trim() ||
                        !formData.riotTag.trim() ||
                        loadingRegions ||
                        securityStatus.suspiciousActivity
                      }
                      sx={{ 
                        mt: 2, 
                        py: 1.5,
                        borderRadius: 2,
                        fontWeight: 600,
                        fontSize: '1.1rem',
                        background: `linear-gradient(135deg, ${leagueColors.blue[500]} 0%, ${leagueColors.blue[600]} 100%)`,
                        '&:hover': {
                          background: `linear-gradient(135deg, ${leagueColors.blue[600]} 0%, ${leagueColors.blue[700]} 100%)`,
                          transform: 'translateY(-1px)',
                          boxShadow: '0 6px 20px rgba(25, 118, 210, 0.3)',
                        },
                        transition: 'all 0.3s ease',
                      }}
                    >
                      {isLoading ? (
                        <>
                          <CircularProgress size={20} sx={{ mr: 1, color: 'inherit' }} />
                          Validation en cours...
                        </>
                      ) : securityStatus.suspiciousActivity ? (
                        "🔒 Compte temporairement verrouillé"
                      ) : (
                        "Valider le compte Riot"
                      )}
                    </Button>
                  </Box>
                </form>
              </CardContent>
            </Card>
          </Grid>

          {/* Google OAuth */}
          <Grid item xs={12}>
            <Card 
              elevation={0}
              sx={{ 
                borderRadius: 3,
                border: `1px solid ${isDarkMode ? leagueColors.dark[200] : leagueColors.blue[100]}`,
                background: isDarkMode
                  ? `linear-gradient(135deg, ${leagueColors.dark[100]} 0%, ${leagueColors.dark[50]} 100%)`
                  : `linear-gradient(135deg, #ffffff 0%, ${leagueColors.gold[25]} 100%)`,
              }}
            >
              <CardContent sx={{ p: 4 }}>
                <Typography 
                  variant="h6" 
                  gutterBottom 
                  sx={{ 
                    textAlign: 'center',
                    mb: 3,
                    color: 'text.primary',
                  }}
                >
                  Connexion alternative
                </Typography>
                <GoogleAuth 
                  onSuccess={(user) => {
                    console.log('OAuth Google réussi:', user);
                  }}
                  onError={(error) => {
                    console.error('Erreur OAuth Google:', error);
                  }}
                />
              </CardContent>
            </Card>
          </Grid>
        </Grid>

        {/* Footer */}
        <Box sx={{ textAlign: "center", mt: 4 }}>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
            Exemple: Riot ID "Canna", Tag "KC", Région "Europe West"
          </Typography>
          <Typography
            variant="caption"
            color="text.secondary"
            sx={{ fontSize: '0.8rem', opacity: 0.8 }}
          >
            🔒 Nous validons votre compte via l'API officielle Riot Games
          </Typography>
        </Box>
      </Box>
    </Box>
  );
}
